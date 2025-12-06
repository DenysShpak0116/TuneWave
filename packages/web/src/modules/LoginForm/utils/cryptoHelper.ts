import JSEncrypt from "jsencrypt";


export function generateClientKeyPair() {
    const crypt = new JSEncrypt({ default_key_size: "2048" });

    const publicKey = crypt.getPublicKey();
    const privateKey = crypt.getPrivateKey()

    return { publicKey, privateKey };
}

export function decryptWithClientPrivateKey(privateKeyPem: string, encryptedBase64: string): string {
    const decryptor = new JSEncrypt();
    decryptor.setPrivateKey(privateKeyPem);

    const decrypted = decryptor.decrypt(encryptedBase64);
    if (!decrypted) {
        throw new Error("Decryption failed");
    }
    return decrypted;
}

export function encryptPassword(serverPublicKeyPem: string, password: string): string {
    const encryptor = new JSEncrypt();
    encryptor.setPublicKey(serverPublicKeyPem);

    const encrypted = encryptor.encrypt(password);
    if (!encrypted) {
        throw new Error("Can't encrypt password");
    }
    return encrypted;
}