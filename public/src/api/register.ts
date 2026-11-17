import { fetch_json } from "./util";
import type { RegisterRequest, RegisterResponse } from "./types";

export const register = (rr: RegisterRequest) => fetch_json<RegisterResponse>(`/api/register`, {
	method: "POST",
	body: JSON.stringify(rr)
});
