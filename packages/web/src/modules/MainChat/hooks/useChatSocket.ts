import { useEffect, useRef, useMemo, useState } from "react"
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
    const participantIds = useMemo(() => usersParam.split(",").filter(Boolean), [usersParam])

    const [qMod, setQMod] = useState<bigint | null>(null)
    const [gGen, setGGen] = useState<bigint | null>(null)
    const [priv, setPriv] = useState<bigint | null>(null)
    const [myPub, setMyPub] = useState<string | null>(null)

    const qModRef = useRef<bigint | null>(null)
    const gGenRef = useRef<bigint | null>(null)
    const privRef = useRef<bigint | null>(null)
    const myPubRef = useRef<string | null>(null)

    const [groupKey, setGroupKey] = useState<CryptoJS.lib.WordArray | null>(null)
    const [groupIv, setGroupIv] = useState<CryptoJS.lib.WordArray | null>(null)
    const seenKeyFp = useRef<Set<string>>(new Set())

    const shouldRegenRef = useRef<boolean>(true)

    useEffect(() => { qModRef.current = qMod }, [qMod])
    useEffect(() => { gGenRef.current = gGen }, [gGen])
    useEffect(() => { privRef.current = priv }, [priv])
    useEffect(() => { myPubRef.current = myPub }, [myPub])

    const regenerateKeys = (socket: WebSocket) => {
        if (!shouldRegenRef.current) return
        console.log("🔁 Перегенерация ключей (новый пользователь)")
        shouldRegenRef.current = false
        seenKeyFp.current.clear()
        setQMod(null); qModRef.current = null
        setGGen(null); gGenRef.current = null
        setPriv(null); privRef.current = null
        setMyPub(null); myPubRef.current = null
        setGroupKey(null); setGroupIv(null)

        const q = randBigInt(256)
        const g = randBigInt(64)
        const a = randBigInt(256) % (q - 3n) + 2n
        const A = modPow(g, a, q)
        const A_b64 = bytesToB64(bigIntToBytes(A))

        setQMod(q); qModRef.current = q
        setGGen(g); gGenRef.current = g
        setPriv(a); privRef.current = a
        setMyPub(A_b64); myPubRef.current = A_b64

        socket.send(JSON.stringify({
            q: bytesToB64(bigIntToBytes(q)),
            p: bytesToB64(bigIntToBytes(g)),
        }))
        for (const uid of participantIds)
            socket.send(JSON.stringify({ receiver: uid, content: A_b64 }))

        setTimeout(() => { shouldRegenRef.current = true }, 2000)
    }

    const tryBuildGroupKeyFromPeerPub = (peerPubB64: string) => {
        const p = qModRef.current
        const g = gGenRef.current
        const a = privRef.current
        if (!p || !a || !g) return
        try {
            seenKeyFp.current.add(peerPubB64)
            const pubs: bigint[] = Array.from(seenKeyFp.current)
                .map((b64) => bytesToBigInt(b64ToBytes(base64urlToBase64(b64))))
                .filter((v) => v > 1n && v < p)
            const myPubNow = myPubRef.current
            if (myPubNow) {
                const myPubBig = bytesToBigInt(b64ToBytes(base64urlToBase64(myPubNow)))
                if (!pubs.includes(myPubBig)) pubs.push(myPubBig)
            }
            pubs.sort((a, b) => (a < b ? -1 : 1))
            let shared = g
            console.log(pubs)
            for (const pub of pubs) shared = modPow(shared, pub, p)
            const sharedB64 = bytesToB64(bigIntToBytes(shared))
            const { key, iv } = deriveAesKeyIv(sharedB64)
            setGroupKey(key)
            setGroupIv(iv)
            localStorage.setItem("message-key", CryptoJS.enc.Base64.stringify(key))
            localStorage.setItem("message-iv", CryptoJS.enc.Base64.stringify(iv))
        } catch { }
    }

    useEffect(() => {
        if (!shouldConnect) return
        const url = `ws://localhost:8081/ws/chat?authToken=${encodeURIComponent(authToken!)}&userIds=${encodeURIComponent(usersParam)}&name=${encodeURIComponent(chatName)}`
        const socket = new WebSocket(url)
        socketRef.current = socket

        socket.onopen = () => console.log("🔗 WS connected:", url)

        socket.onmessage = (event) => {
            try {
                let data: any
                try { data = JSON.parse(event.data) } catch { data = event.data }

                if (data?.q && data?.p) {
                    const q = bytesToBigInt(b64ToBytes(base64urlToBase64(data.q)))
                    const g = bytesToBigInt(b64ToBytes(base64urlToBase64(data.p)))
                    setQMod(q); qModRef.current = q
                    setGGen(g); gGenRef.current = g
                    const a = randBigInt(256) % (q - 3n) + 2n
                    const A = modPow(g, a, q)
                    const A_b64 = bytesToB64(bigIntToBytes(A))
                    setPriv(a); privRef.current = a
                    setMyPub(A_b64); myPubRef.current = A_b64
                    console.log("🧩 Мой публичный ключ:", A_b64)
                    for (const uid of participantIds)
                        socket.send(JSON.stringify({ receiver: uid, content: A_b64 }))
                    return
                }

                if (typeof data === "string" && isProbablyB64(data)) {
                    const pubB64 = data.trim()
                    const myPubNow = myPubRef.current
                    if (myPubNow && pubB64 === myPubNow) return

                    if (!seenKeyFp.current.has(pubB64)) {
                        seenKeyFp.current.add(pubB64)
                        console.log("📥 Получен чужой публичный ключ")

                        if (shouldRegenRef.current) {
                            regenerateKeys(socket)
                            return
                        }

                        if (myPubNow)
                            for (const uid of participantIds)
                                socket.send(JSON.stringify({ receiver: uid, content: myPubNow }))
                        tryBuildGroupKeyFromPeerPub(pubB64)
                    }
                    return
                }

                if (typeof data === "object" && typeof data.content === "string") {
                    if (groupKey) {
                        const plain = aesDecrypt(data.content, groupKey, groupIv!)
                        if (plain) {
                            onMessageReceived({ ...data, content: plain })
                            return
                        }
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
            try { socket.close() } catch { }
            socketRef.current = null
            seenKeyFp.current.clear()
            setQMod(null); qModRef.current = null
            setPriv(null); privRef.current = null
            setMyPub(null); myPubRef.current = null
            setGroupKey(null); setGroupIv(null)
        }
    }, [shouldConnect, authToken, usersParam])

    const sendMessage = (content: string) => {
        const socket = socketRef.current
        if (!socket || socket.readyState !== WebSocket.OPEN) {
            console.warn("⛔ WS не открыт")
            return
        }
        if (groupKey && groupIv) {
            const payload = aesEncrypt(content, groupKey, groupIv)
            socket.send(JSON.stringify({ content: payload }))
        } else {
            console.warn("⚠️ Общий ключ не готов — сообщение отправлено открытым")
            socket.send(JSON.stringify({ content }))
        }
    }

    return { sendMessage }
}
