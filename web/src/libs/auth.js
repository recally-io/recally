import Cookies from "js-cookie";
import { get } from "./api";

export async function checkIsLogin() {
  const token = Cookies.get("token");
  if (!token) {
    return false;
  }
  try {
    const res = await get("/api/v1/auth/validate-jwt");
    console.debug("Token is valid, user is logged in", res.data);
    return true;
  } catch (error) {
    console.debug("Token is invalid", error);
    return false;
  }
}
