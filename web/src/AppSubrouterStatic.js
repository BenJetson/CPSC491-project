import React from "react";
import { Route, Switch, useRouteMatch } from "react-router-dom";
import NotFound from "./components/NotFound";
import FAQ from "./components/FAQ";
import About from "./components/About";
const AppSubrouterStatic = () => {
  const match = useRouteMatch();

  return (
    <Switch>
      <Route path={`${match.path}/help`}>
        <FAQ />
      </Route>
      <Route path={`${match.path}/about`}>
        <About />
      </Route>
      <Route path={"*"}>
        {/* If no route matches, show a not found page. */}
        <NotFound />
      </Route>
    </Switch>
  );
};

export default AppSubrouterStatic;
