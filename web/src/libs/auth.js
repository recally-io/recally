import Cookies from "js-cookie";
import { authApi } from "./api";

export async function checkIsLogin() {
  const token = Cookies.get("token");
  if (!token) {
    console.log("No token found, redirecting to login page");
    return false;
  }
  try {
    await authApi.authValidateJwtGet({ token: token });
    console.debug("Token is valid");
    return true;
  } catch (error) {
    console.debug("Token is invalid", error);
    return false;
  }
}
