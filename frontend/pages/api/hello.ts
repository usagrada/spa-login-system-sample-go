// Next.js API route support: https://nextjs.org/docs/api-routes/introduction
import type { NextApiRequest, NextApiResponse } from 'next'
import axios from "axios"
axios.defaults.xsrfCookieName = "_csrf"
axios.defaults.xsrfHeaderName = "X-CSRF-Token"
axios.defaults.withCredentials = true

type Data = {
  name: string
}

export default async function handler(
  req: NextApiRequest,
  res: NextApiResponse<Data>
) {
  const csrf_res = await axios.get("http://localhost:8080/api/csrf")

  const login_res = await axios.post("http://localhost:8080/api/signup", {
    username: "admin",
    password: "admin",
    // csrfToken: csrf
  })
  console.log(login_res.data)
  res.status(200).json({ name: 'John Doe' })
}
