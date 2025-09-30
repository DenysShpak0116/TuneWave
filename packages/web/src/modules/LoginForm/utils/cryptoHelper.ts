import JSEncrypt from "jsencrypt";

function arrayBufferToBase64(buffer: ArrayBuffer): string {
    let binary = "";
    const bytes = new Uint8Array(buffer);
    for (let i = 0; i < bytes.byteLength; i++) {
        binary += String.fromCharCode(bytes[i]);
    }
    return window.btoa(binary);
}

export async function generateClientKeyPair() {
    const keyPair = await window.crypto.subtle.generateKey(
        {
            name: "RSA-OAEP",
            modulusLength: 2048,
            publicExponent: new Uint8Array([1, 0, 1]),
            hash: "SHA-256",
        },
        true,
        ["encrypt", "decrypt"]
    );

    const publicKey = await window.crypto.subtle.exportKey("spki", keyPair.publicKey);
    const privateKey = await window.crypto.subtle.exportKey("pkcs8", keyPair.privateKey);

    return {
        publicKey: arrayBufferToBase64(publicKey),
        privateKey: arrayBufferToBase64(privateKey),
    };
}

// function pemToBase64(pem: string): string {
//     return pem
//         .replace(/-----BEGIN PUBLIC KEY-----/, "")
//         .replace(/-----END PUBLIC KEY-----/, "")
//         .replace(/\s+/g, "");
// }

// function base64ToArrayBuffer(base64: string): ArrayBuffer {
//     const binary = window.atob(base64);
//     const bytes = new Uint8Array(binary.length);
//     for (let i = 0; i < binary.length; i++) {
//         bytes[i] = binary.charCodeAt(i);
//     }
//     return bytes.buffer;
// }

export function encryptPassword(serverPublicKeyPem: string, password: string): string {
    const encryptor = new JSEncrypt();
    encryptor.setPublicKey(serverPublicKeyPem);

    const encrypted = encryptor.encrypt(password);
    if (!encrypted) {
        throw new Error("Не удалось зашифровать пароль");
    }
    return encrypted;
}