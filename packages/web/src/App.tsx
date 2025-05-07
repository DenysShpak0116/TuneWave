import GlobalStyle from "./global-style"
import { Toaster } from "react-hot-toast"
import { COLORS } from "@consts/colors.consts"
import { FONTS } from "@consts/fonts.enum"
import { checkAuth } from "helpers/Auth/check-auth"
import { useEffect } from "react"
import { AppRoutes } from "pages/router/app-routes"
import "./helpers/Player/PlayerSync";
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
      <AppRoutes />
    </>
  )
}

export default App
