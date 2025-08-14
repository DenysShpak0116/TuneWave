export const parseDate = (date?: string): string => {
    if (!date || !date.includes('T')) return "Дата не вказана";
    return date.split('T')[0];
};

export const parseTime = (timeStr: string): string => {
    const regex = /(?:(\d+)m)?([\d.]+)s/;
    const match = timeStr.match(regex);

    if (!match) {
        throw new Error('Не правильний формат часу');
    }
    const minutes = match[1] ? parseInt(match[1], 10) : 0;
    const secondsFloat = parseFloat(match[2]);
    const totalSeconds = Math.round(minutes * 60 + secondsFloat);

    const minutesPart = Math.floor((totalSeconds % 3600) / 60);
    const secondsPart = totalSeconds % 60;

    const MM = minutesPart.toString().padStart(2, '0');
    const SS = secondsPart.toString().padStart(2, '0');

    return `${MM}:${SS}`;
}