import Cookies from "js-cookie";
import { AuthApi } from "../sdk/index";

export async function checkIsLogin() {
  const token = Cookies.get("token");
  if (!token) {
    console.log("No token found, redirecting to login page");
    return false;
  }
  console.log("Checking token validity");
  try {
    await new AuthApi().authValidateJwtGet({ token: token });
    console.log("Token is valid");
    return true;
  } catch (error) {
    console.debug("Token is invalid", error);
    return false;
  }
}
