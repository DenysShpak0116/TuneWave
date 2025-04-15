import { FONTS } from "consts/fonts.enum"
import GlobalStyle from "./global-style"

function App() {

  return (
    <>
      <GlobalStyle />
      <h1 style={{ fontFamily: `${FONTS.MONTSERRAT}` }}>Hello page!</h1>
    </>
  )
}

export default App
