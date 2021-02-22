import { createMuiTheme } from "@material-ui/core/styles";

//  The default theme of the webpage
//  https://material-ui.com/customization/color/#color
//  for info on how color schemes are generated
const defaultTheme = createMuiTheme({
    palette: {
      primary: {
        main: '#3f50b5',
      },
      secondary: {
        main: '#f44336',
      },
    },
  });

  export default defaultTheme;