import { OperationResult } from "../../../data/repositories/_types/types";
import { AuthRepo } from "../../../infra/respositories/local-storage/auth/auth-impl.repo";


export async function getAccessTokens(): Promise<OperationResult<string>> {

    const repo = new AuthRepo();

    const getTokensRes = await repo.getTokens();
    if (!getTokensRes || !getTokensRes.accessToken) 
        return {
            success: false,
            message: 'No access token available'
        }

    return { 
        success: true, 
        data: getTokensRes.accessToken 
    };
}