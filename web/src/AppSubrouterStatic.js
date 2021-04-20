import React from "react";
import { Route, Switch, useRouteMatch } from "react-router-dom";
import NotFound from "./components/NotFound";
import NotImplemented from "./components/NotImplemented";

const AppSubrouterStatic = () => {
  const match = useRouteMatch();

  return (
    <Switch>
      <Route path={`${match.path}/help`}>
        <NotImplemented feature={"Static - Help"} />
      </Route>
      <Route path={`${match.path}/about`}>
        <NotImplemented feature={"Static - About"} />
      </Route>
      <Route path={"*"}>
        {/* If no route matches, show a not found page. */}
        <NotFound />
      </Route>
    </Switch>
  );
};

export default AppSubrouterStatic;
