import GlobalStyle from "./global-style"
import { RouterProvider } from "react-router-dom"
import router from "pages/router"

function App() {

  return (
    <>
      <GlobalStyle />
      <RouterProvider router={router} />
    </>
  )
}

export default App
