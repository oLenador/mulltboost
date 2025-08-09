import { UserSettingsRepo } from "../../../infra/respositories/local-storage/user-impl.repo";


export function GetUserSettingsUseCase() {
    let userSettingsRepo: UserSettingsRepo = new UserSettingsRepo();
    return userSettingsRepo.get();
}