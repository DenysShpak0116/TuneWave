import CryptoJS from "crypto-js";

export const deriveAesKeyIv = (sharedB64: string) => {
    const hash = CryptoJS.SHA256(CryptoJS.enc.Utf8.parse(sharedB64));
    const hex = CryptoJS.enc.Hex.stringify(hash);
    const key = CryptoJS.enc.Hex.parse(hex);
    const iv = CryptoJS.enc.Hex.parse(hex.slice(0, 32));
    return { key, iv };
};

export const aesEncrypt = (
    plaintext: string,
    key: CryptoJS.lib.WordArray,
    iv: CryptoJS.lib.WordArray
): string => {
    return CryptoJS.AES.encrypt(plaintext, key, {
        iv, mode: CryptoJS.mode.CBC, padding: CryptoJS.pad.Pkcs7,
    }).toString();
};

export const aesDecrypt = (
    cipherB64: string,
    key: CryptoJS.lib.WordArray,
    iv: CryptoJS.lib.WordArray
): string | null => {
    try {
        const pt = CryptoJS.AES.decrypt(
            // eslint-disable-next-line @typescript-eslint/no-explicit-any
            { ciphertext: CryptoJS.enc.Base64.parse(cipherB64) } as any,
            key,
            { iv, mode: CryptoJS.mode.CBC, padding: CryptoJS.pad.Pkcs7 }
        );
        return CryptoJS.enc.Utf8.stringify(pt);
    } catch {
        return null;
    }
};