export const parseDate = (date?: string): string => {
    if (!date || !date.includes('T')) return "Дата не вказана";
    return date.split('T')[0];
};