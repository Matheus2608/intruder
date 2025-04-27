import { ResponseData } from "./ResponseData";

export interface AttackOutput {
    responses: ResponseData[],
    url: string,
    totalElapsedTimeInMiliseconds : number
}
