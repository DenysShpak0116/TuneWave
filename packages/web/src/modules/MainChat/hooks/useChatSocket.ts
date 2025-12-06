import { useEffect, useRef, useState } from "react"
import CryptoJS from "crypto-js"
import { IMessageType } from "types/chat/message.type"
import { deriveAesKeyIv, aesEncrypt, aesDecrypt } from "./cryptoAes"

const base64urlToBase64 = (s: string): string =>
    (s || "").replace(/-/g, "+").replace(/_/g, "/") + "=".repeat((4 - ((s || "").length % 4)) % 4)

const b64ToBytes = (b64: string): Uint8Array => {
    const fixed = (b64 || "").replace(/-/g, "+").replace(/_/g, "/")
    const pad = fixed + "=".repeat((4 - fixed.length % 4) % 4)
    const bin = atob(pad)
    const out = new Uint8Array(bin.length)
    for (let i = 0; i < bin.length; i++) out[i] = bin.charCodeAt(i)
    return out
}

const bytesToB64 = (bytes: Uint8Array): string => {
    let bin = ""
    for (let i = 0; i < bytes.length; i++) bin += String.fromCharCode(bytes[i])
    return btoa(bin)
}

const bytesToBigInt = (bytes: Uint8Array): bigint =>
    bytes.reduce((v, b) => (v << 8n) | BigInt(b), 0n)

const bigIntToBytes = (v: bigint): Uint8Array => {
    if (v === 0n) return new Uint8Array([0])
    const buf: number[] = []
    for (let x = v; x > 0n; x >>= 8n) buf.push(Number(x & 0xffn))
    buf.reverse()
    return new Uint8Array(buf)
}

const modPow = (base: bigint, exp: bigint, mod: bigint): bigint => {
    let r = 1n, b = base % mod, e = exp
    while (e > 0n) {
        if (e & 1n) r = (r * b) % mod
        b = (b * b) % mod
        e >>= 1n
    }
    return r
}

const randBigInt = (bits: number): bigint => {
    const len = Math.ceil(bits / 8)
    const buf = new Uint8Array(len)
    crypto.getRandomValues(buf)
    const extra = len * 8 - bits
    if (extra > 0) buf[0] &= (1 << (8 - extra)) - 1
    let v = bytesToBigInt(buf)
    if (v < 2n) v = 2n
    return v
}

const isProbablyB64 = (s: string) => /^[A-Za-z0-9+/_-]+={0,2}$/.test(s) && s.length >= 16

interface UseChatSocketProps {
    usersParam: string
    chatName: string
    onMessageReceived: (message: IMessageType | { __reset: true }) => void
}

export const useChatSocket = ({
    usersParam,
    chatName,
    onMessageReceived,
}: UseChatSocketProps) => {
    const socketRef = useRef<WebSocket | null>(null)
    const authToken = localStorage.getItem("token")
    const shouldConnect = !!authToken && usersParam.trim() !== "" && chatName.trim() !== ""

    const participantIdsRef = useRef<string[]>([])
    useEffect(() => {
        participantIdsRef.current = usersParam.split(",").filter(Boolean)
    }, [usersParam])

    const qModRef = useRef<bigint | null>(null)
    const gGenRef = useRef<bigint | null>(null)
    const privRef = useRef<bigint | null>(null)
    const myPubRef = useRef<string | null>(null)

    const [groupKey, setGroupKey] = useState<CryptoJS.lib.WordArray | null>(null)
    const [groupIv, setGroupIv] = useState<CryptoJS.lib.WordArray | null>(null)
    // eslint-disable-next-line @typescript-eslint/no-unused-vars
    const [decryptedMessages, setDecryptedMessages] = useState<IMessageType[]>([])

    const seenKeyFp = useRef<Set<string>>(new Set())

    const tryBuildGroupKey = () => {
        const p = qModRef.current
        const g = gGenRef.current
        const a = privRef.current
        const myPub = myPubRef.current

        if (!p || !a || !g || !myPub) return

        if (!seenKeyFp.current.has(myPub)) {
            seenKeyFp.current.add(myPub)
        }

        try {
            const pubs: bigint[] = Array.from(seenKeyFp.current)
                .map((b64) => bytesToBigInt(b64ToBytes(base64urlToBase64(b64))))
                .filter((v) => v > 1n && v < p)

            pubs.sort((a, b) => (a < b ? -1 : 1))

            let shared = g

            console.log(pubs)
            for (const pub of pubs) {
                shared = modPow(shared, pub, p)
            }

            const sharedB64 = bytesToB64(bigIntToBytes(shared))
            const { key, iv } = deriveAesKeyIv(sharedB64)

            setGroupKey(key)
            setGroupIv(iv)

            localStorage.setItem("message-key", CryptoJS.enc.Base64.stringify(key))
            localStorage.setItem("message-iv", CryptoJS.enc.Base64.stringify(iv))
        } catch (e) {
            console.error(e)
        }
    }


    const resetAndStartDhExchange = (q: bigint, g: bigint, socket: WebSocket) => {
        seenKeyFp.current.clear()
        setGroupKey(null); setGroupIv(null)

        const a = randBigInt(256) % (q - 3n) + 2n
        const A = modPow(g, a, q)
        const A_b64 = bytesToB64(bigIntToBytes(A))

        qModRef.current = q
        gGenRef.current = g
        privRef.current = a
        myPubRef.current = A_b64

        seenKeyFp.current.add(A_b64)

        for (const uid of participantIdsRef.current) {
            socket.send(JSON.stringify({ receiver: uid, content: A_b64 }))
        }

        tryBuildGroupKey()
    }

    useEffect(() => {
        if (!shouldConnect) return

        const url = `ws://localhost:8081/ws/chat?authToken=${encodeURIComponent(authToken!)}&userIds=${encodeURIComponent(usersParam)}&name=${encodeURIComponent(chatName)}`
        const socket = new WebSocket(url)
        socketRef.current = socket

        socket.onopen = () => console.log("🔗 WS connected:", url)

        socket.onmessage = (event) => {
            try {
                // eslint-disable-next-line @typescript-eslint/no-explicit-any
                let data: any
                try { data = JSON.parse(event.data) } catch { data = event.data }
                console.log(data)

                if (data?.q && data?.p) {
                    const modulus = bytesToBigInt(b64ToBytes(base64urlToBase64(data.p)))
                    let generator = bytesToBigInt(b64ToBytes(base64urlToBase64(data.q)))

                    generator %= modulus
                    if (generator <= 1n) generator = 2n

                    resetAndStartDhExchange(modulus, generator, socket)
                    return
                }

                if (typeof data === "string" && isProbablyB64(data)) {
                    const pubB64 = data.trim()
                    const myPubNow = myPubRef.current

                    if (myPubNow && pubB64 === myPubNow) return

                    if (!seenKeyFp.current.has(pubB64)) {
                        seenKeyFp.current.add(pubB64)

                        if (myPubNow) {
                            for (const uid of participantIdsRef.current)
                                socket.send(JSON.stringify({ receiver: uid, content: myPubNow }))
                        }

                        tryBuildGroupKey()
                    }
                    return
                }

                if (typeof data === "object" && typeof data.content === "string") {
                    let plain: string | null = null

                    if (groupKey && groupIv) {
                        plain = aesDecrypt(data.content, groupKey, groupIv!)
                    }

                    if (!plain) {
                        const prevKeyB64 = localStorage.getItem("message-key")
                        const prevIvB64 = localStorage.getItem("message-iv")
                        if (prevKeyB64 && prevIvB64) {
                            try {
                                const prevKey = CryptoJS.enc.Base64.parse(prevKeyB64)
                                const prevIv = CryptoJS.enc.Base64.parse(prevIvB64)
                                plain = aesDecrypt(data.content, prevKey, prevIv)
                            } catch (err) {
                                console.warn(err)
                            }
                        }
                    }

                    if (plain) {
                        const msg: IMessageType = { ...data, content: plain }

                        setDecryptedMessages(prev => {
                            const updated = [...prev, msg]
                            localStorage.setItem("chat-decrypted", JSON.stringify(updated))
                            return updated
                        })

                        onMessageReceived(msg)
                        return
                    }

                    const saved = JSON.parse(localStorage.getItem("chat-decrypted") || "[]")
                    const found = saved.find((m: IMessageType) => m.id === data.id)
                    if (found) {
                        onMessageReceived(found)
                        return
                    }

                    onMessageReceived(data as IMessageType)
                }
            } catch (e) {
                console.error("❌ Failed to handle WS message:", e)
            }
        }

        socket.onclose = () => console.log("❌ WS closed")
        socket.onerror = (e) => console.error("⚠️ WS error:", e)

        return () => {
            try { socket.close() } catch {
                console.error("error while closing socket connection")
            }
            socketRef.current = null
            // eslint-disable-next-line react-hooks/exhaustive-deps
            seenKeyFp.current.clear()
            qModRef.current = null
            gGenRef.current = null
            privRef.current = null
            myPubRef.current = null
            setGroupKey(null); setGroupIv(null)
        }
    // eslint-disable-next-line react-hooks/exhaustive-deps
    }, [shouldConnect, authToken, usersParam])

    const sendMessage = (content: string) => {
        const socket = socketRef.current
        if (!socket || socket.readyState !== WebSocket.OPEN) {
            return
        }
        if (groupKey && groupIv) {
            const payload = aesEncrypt(content, groupKey, groupIv)
            socket.send(JSON.stringify({ content: payload }))
        } else {
            socket.send(JSON.stringify({ content }))
        }
    }

    return { sendMessage }
}
