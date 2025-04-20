import { T } from "../wailsjs/go/main/App"

export default async (message) => {
  return await T(message, navigator.language)
}
