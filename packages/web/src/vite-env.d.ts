/// <reference types="vite/client" />

declare module '*.ttf' {
    const content: string
    export default content
}

declare global {
    interface Window {
        PlayerjsEvents: (event: string, id: string, info: any) => void;
    }
}
export { };