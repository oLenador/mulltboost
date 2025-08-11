import { z } from "zod";
import { OperationResult } from "../../../data/repositories/_types/types";
import { JwtAuthPayload } from "../../../data/repositories/i-authentication.repo";
import { JwtPayloadSchema } from "../../validation/jwt-payload.schema";
import { jwtDecode } from "jwt-decode";
import { getAccessTokens } from "./get-access-token.usecase";


export async function getJwtPayload(): Promise<OperationResult<JwtAuthPayload>> {
    // console.log("TeTent!")

    const getAccessTokenRes = await getAccessTokens();
    // console.log(getAccessTokenRes)

    if (!getAccessTokenRes.success || !getAccessTokenRes.data) return {
        success: false,
        message: 'No access token available'
    };

    const decodedPayloaRes = await decodeJwtPayload(getAccessTokenRes.data, JwtPayloadSchema);
    if (!decodedPayloaRes.success || !decodedPayloaRes.data) return {
        success: false,
        message: 'No access token available'
    };
    // console.log("TenTentTentTent!")

    const getJwtPayloadRes = decodedPayloaRes.data;
    return { 
        success: true, 
        data: getJwtPayloadRes 
    }
}


export async function decodeJwtPayload(
    jwtPayload: string, 
    schemaValidation: z.ZodObject<any>
): Promise<OperationResult<JwtAuthPayload>> {


    try {
        const decodePayloadRes = jwtDecode(jwtPayload);
        if (!decodePayloadRes || typeof decodePayloadRes !== 'object') return {
            success: false,
            message: 'No access token available'
        };        
        // console.log(decodePayloadRes)
        // console.log("validating schema")
        // schemaValidation.parse(decodePayloadRes)
        // console.log("dwadwa schema")

        const decodedRes = decodePayloadRes as JwtAuthPayload

        return {
            success: true,
            data: decodedRes
        }
    } catch (error) {
        // console.log(error)
        return {
            success: false,
            message: error.message
        }
    }
}