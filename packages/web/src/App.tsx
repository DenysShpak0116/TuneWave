import GlobalStyle from "./global-style"
import { RouterProvider } from "react-router-dom"
import router from "pages/router"
import { Toaster } from "react-hot-toast"
import { COLORS } from "@consts/colors.consts"
import { FONTS } from "@consts/fonts.enum"
import { checkAuth } from "helpers/Auth/check-auth"
import { useEffect } from "react"

function App() {

  useEffect(() => {
    checkAuth();
  }, []);

  return (
    <>
      <Toaster
        position="top-right"
        toastOptions={{
          style: {
            background: `${COLORS.dark_focusing}`,
            color: `${COLORS.dark_secondary}`,
            fontFamily: `${FONTS.MONTSERRAT}`,
            fontSize: "16px",
            borderRadius: "5px",
          },
        }}
      />
      <GlobalStyle />
      <RouterProvider router={router} />
    </>
  )
}

export default App
