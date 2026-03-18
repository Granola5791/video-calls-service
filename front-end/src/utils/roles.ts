import { Roles } from "../constants/general-contants";

export const IsAdmin = (role: string | null): boolean => {
    return role !== null && role.includes(Roles.admin);
}