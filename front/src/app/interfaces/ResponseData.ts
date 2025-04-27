export interface ResponseData {
    requestId: number,
    payload: string,
    statusCode: number,
    elapsedTimeInMilliseconds: number,
    didError: boolean,
    bytesLenght: number,
    request: string,
    response: string
}
