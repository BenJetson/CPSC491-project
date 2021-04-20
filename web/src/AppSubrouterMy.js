import React from "react";
import { Route, Switch, useRouteMatch } from "react-router-dom";
import NotFound from "./components/NotFound";
import NotImplemented from "./components/NotImplemented";

const AppSubrouterMy = () => {
  const match = useRouteMatch();

  return (
    <Switch>
      <Route path={`${match.path}/profile`}>
        <NotImplemented feature={"My - Profile"} />
      </Route>
      <Route path={"*"}>
        {/* If no route matches, show a not found page. */}
        <NotFound />
      </Route>
    </Switch>
  );
};

export default AppSubrouterMy;
