import { createContext, useEffect, useState } from "react";
import { JwtAuthPayload } from "../../data/repositories/i-authentication.repo";
import { getJwtPayload } from "../../domain/usecases/auth/get-jwt-payload";

interface MineAccount {

}

interface Servers {
    title: string
    ip: string
}

interface UserProviderI {
    accounts: Array<MineAccount>
    servers: Array<Servers>
}

export const UserProvider_INITIAL: UserProviderI = {
    accounts: [],
    servers: []
}
export const UserProvider = createContext<UserProviderI>(UserProvider_INITIAL)



export function UserProviderHook(): UserProviderI {
    const [user, setUser] = useState<UserProviderI>(UserProvider_INITIAL);
    console.log(user)
    async function loadData() {
        // const checkLoginRes = await CheckLoginUseCase()
// 
        // if (!checkLoginRes.success) {
        //     setUser(UserProvider_INITIAL)
        //     return;
        // }
// 
        // const payload = await getJwtPayload();
        // if (!payload.success || !payload.data) {
        //     setUser(UserProvider_INITIAL)
        //     return;
        // }
        // setUser({
        //     authenticated,
        //     ...payload.data
        // })
    }

    useEffect(() => {
        loadData();
    }, [])


    return user;
}