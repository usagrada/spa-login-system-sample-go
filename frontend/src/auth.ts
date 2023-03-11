import axios from "axios";
import { NextApiResponse } from "next";

axios.defaults.xsrfCookieName = "_csrf";
axios.defaults.xsrfHeaderName = "X-CSRF-Token";
axios.defaults.withCredentials = true;

export async function getJwt(res: NextApiResponse<any>) {
  const jwt_res = await axios.get("http://localhost:8080/jwt");
  console.log(jwt_res.data)
  axios.defaults.headers.common["Authorization"] = `Bearer ${jwt_res.data.token}`;
  res.setHeader("Set-Cookie", jwt_res.headers["set-cookie"] ?? []);
}
