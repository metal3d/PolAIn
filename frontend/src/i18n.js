import { T } from "../wailsjs/go/main/App"

export default async (message, md = false) => {
  return await T(message, navigator.language, md)
}
