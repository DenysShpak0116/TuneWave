import { ICompareType } from "./compare.type";

export interface IResultType {
    song1Id: string;
    comparedTo: ICompareType[]
}