import Cookies from "js-cookie";
import { request } from "./api";

export async function checkIsLogin() {
  const token = Cookies.get("token");
  if (!token) {
    return false;
  }
  try {
    const res = await request("/api/v1/auth/validate-jwt");
    const data = res.json();
    console.debug("Token is valid, user is logged in", data.data);
    return true;
  } catch (error) {
    console.debug("Token is invalid", error);
    return false;
  }
}
