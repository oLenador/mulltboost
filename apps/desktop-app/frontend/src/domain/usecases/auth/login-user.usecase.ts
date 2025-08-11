import { OperationResult } from "../../../data/repositories/_types/types";
import { AuthRepo } from "../../../infra/respositories/local-storage/auth/auth-impl.repo";

export function LoginUserUseCase(token: string): Promise<OperationResult> {
    let authRepo: AuthRepo = new AuthRepo();
    return authRepo.loginUser(token);
}

