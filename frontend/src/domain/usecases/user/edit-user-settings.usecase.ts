import { OperationResult } from "../../../data/repositories/_types/types";
import { UserSettingsRepo } from "../../../infra/respositories/local-storage/user-impl.repo";
import { UserSettingsI } from "../../entity/user.entity";



export async function EditUserSettingsUseCase(settings: UserSettingsI): Promise<OperationResult> {
    let userSettingsRepo: UserSettingsRepo = new UserSettingsRepo();
    return userSettingsRepo.update(settings);
}