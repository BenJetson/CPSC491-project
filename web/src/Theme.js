import { createMuiTheme } from "@material-ui/core/styles";
const mainColor = "#a01a58";
const secColor = "#1780a1";

//  The default theme of the webpage
//  https://material-ui.com/customization/color/#color
//  for info on how color schemes are generated
const defaultTheme = createMuiTheme({
  palette: {
    type: "light",
    primary: {
      main: mainColor,
    },
    secondary: {
      main: secColor,
    },
  },
});

export default defaultTheme;
