import { OperationResult } from "../../../data/repositories/_types/types";
import { JwtAuthPayload } from "../../../data/repositories/i-authentication.repo";
import { getJwtPayload } from "./get-jwt-payload";

export async function CheckLoginUseCase(): Promise<OperationResult> {
// let authRepo: AuthRepo = new AuthRepo();
// // console.log("---------------Checking login---------------")
// const validToken = await verifyTokenValidation()
// if (validToken.success && validToken.data) return {
//     success: true
// }
// // console.log(validToken)
// const valideAccesToken = await authRepo.validateAccessToken();
// if (!valideAccesToken.success) return {
//     success: false,
//     message: 'Token inválido.'
// }
// 
 return {
     success: true
 }
}
