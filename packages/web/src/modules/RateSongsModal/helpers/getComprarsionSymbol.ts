export const getComparisonSymbol = (value: number) => {
    switch (value) {
        case -2: return "<";
        case -1: return "<=";
        case 0: return "=";
        case 1: return ">=";
        case 2: return ">";
        default: return "";
    }
};