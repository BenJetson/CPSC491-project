import React from "react";
import { Route, Switch, useRouteMatch } from "react-router-dom";

import NotImplemented from "./components/NotImplemented";
import NotFound from "./components/NotFound";

const AppSubrouterSponsor = () => {
  const match = useRouteMatch();

  return (
    <Switch>
      <Route path={`${match.path}/applications`}>
        <NotImplemented feature={"Sponsor - Manage Applications"} />
      </Route>
      <Route path={`${match.path}/drivers`}>
        <NotImplemented feature={"Sponsor - Manage Drivers"} />
      </Route>
      <Route path={`${match.path}/catalog`}>
        <NotImplemented feature={"Sponsor - Manage Catalog"} />
      </Route>
      <Route path={"*"}>
        {/* If no route matches, show a not found page. */}
        <NotFound />
      </Route>
    </Switch>
  );
};

export default AppSubrouterSponsor;
