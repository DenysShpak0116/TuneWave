type Song = {
    duration: string;
};

export const getTotalDuration = (songs: Song[]): string => {
    let totalSeconds = 0;

    songs.forEach((song) => {
        const [minutes, seconds] = song.duration.split(":").map(Number);
        totalSeconds += minutes * 60 + seconds;
    });

    const hours = Math.floor(totalSeconds / 3600);
    const minutes = Math.floor((totalSeconds % 3600) / 60);
    const seconds = totalSeconds % 60;

    const format = (n: number) => String(n).padStart(2, "0");

    return `${hours > 0 ? format(hours) + ":" : ""}${format(minutes)}:${format(seconds)}`;
};