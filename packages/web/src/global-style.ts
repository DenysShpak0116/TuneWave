import { createGlobalStyle } from "styled-components";
import { FONTS } from "@consts/fonts.enum";
import { COLORS } from "@consts/colors.consts";


const GlobalStyle = createGlobalStyle`
* {
      margin: 0;
      padding: 0;
      font-family: '${FONTS.MONTSERRAT}';
  }

  body{
    background-color: ${COLORS.dark_backdrop};
    padding-bottom: 50px;
  }

  h1,
  h2,
  h3,
  h4,
  h5,
  h6 {
      margin: 0;
  }

  a {
      display: block;
      text-decoration: none;
  }

  ul,
  ol {
      margin: 0;
      padding: 0;
      list-style: none;
  }

  button {
      cursor: pointer;
      border: none;
      outline: none !important;
  }

  input {
      outline: none !important;
      border: none;
  }

  @font-face {
    font-family: ${FONTS.MONTSERRAT};
    src: url('/fonts/Montserrat.ttf') format('truetype');
    ont-display: swap;
  }

  @font-face {
    font-family: ${FONTS.JERSEY25};
    src: url('/fonts/Jersey25-Regular.ttf') format('truetype');
    font-display: swap;
  }

`
export default GlobalStyle