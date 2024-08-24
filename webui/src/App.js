// PPH MODIFIED
import React, { useReducer } from "react";
import * as R from "ramda";
import styled from 'styled-components';
import { ThemeProvider } from '@mui/material/styles';
import CssBaseline from '@mui/material/CssBaseline';
import { StyledEngineProvider } from '@mui/material/styles';
import Form from "./Form";
import Results from "./Results";
import { Container } from "./Components";
import { buildQueryString, isNilOrEmpty } from "./utils";
import { search } from "./api";
import { createBrowserHistory } from "history";
import * as jose from "jose";
import Cookies from "js-cookie";
import theme from "./theme";
import StylesProvider from "./StylesProvider";

const history = createBrowserHistory();

const MarginDiv = styled.div`
  margin: 1em auto;
`;

const StyledLink = styled.a`
  color: #0000ee;
`;

const createJWT = async (apiKey) =>
{
  const secret = new TextEncoder().encode(apiKey);
  const alg = "HS256";

  const time = new Date();
  time.setMinutes(time.getMinutes() - 1);
  const timeNum = time.getTime() / 1000;

  const jwt = await new jose.SignJWT({ "urn:example:claim": true })
    .setProtectedHeader({ alg })
    .setIssuedAt(timeNum)
    .setIssuer("urn:example:issuer")
    .setAudience("urn:example:audience")
    .setExpirationTime("24h")
    .sign(secret);

  console.log("JWT: " + jwt);

  return jwt;
};

const reducer = (state, action) =>
{
  switch (action.type)
  {
    case "SEARCH_INIT":
      return R.pipe(
        R.assoc("results", null),
        R.assoc("loading", true),
        R.assoc("error", null)
      )(state);
    case "SEARCH_SUCCESS":
      return R.pipe(
        R.assoc("results", action.payload),
        R.assoc("loading", false),
        R.assoc("error", null)
      )(state);
    case "SEARCH_ERROR":
      return R.pipe(
        R.assoc("results", null),
        R.assoc("loading", false),
        R.assoc("error", action.payload)
      )(state);
    case "SEARCH_RESET":
      return initialState;
    default:
      return state;
  }
};

const initialState = {
  error: null,
  loading: false,
  results: null,
};

const valuesOnlyContainLimit = R.pipe(
  R.filter(R.complement(isNilOrEmpty)),
  R.omit(["limit"]),
  R.isEmpty
);

function App()
{
  const [state, dispatch] = useReducer(reducer, initialState);
  const executeSearch = async (qs, apiKey) =>
  {
    dispatch({ type: "SEARCH_INIT" });
    try
    {
      const payload = await search(qs, apiKey);
      dispatch({ type: "SEARCH_SUCCESS", payload });
    } catch (err)
    {
      dispatch({ type: "SEARCH_ERROR", payload: err });
    }
  };

  const handleReset = () =>
  {
    dispatch({ type: "SEARCH_RESET" });
    history.push({ ...history.location, search: "" });
  };

  const handleSubmit = async (values) =>
  {
    if (valuesOnlyContainLimit(values)) return;

    const apiKey = values.apiKey;
    const token = await createJWT(apiKey);

    Cookies.set("token", token, { expires: 1 });

    const valuesWithoutApiKey = R.omit(["apiKey"], values);

    const qs = buildQueryString(valuesWithoutApiKey);
    history.push({ ...history.location, search: qs });
    executeSearch(qs);
  };

  return (
    <StylesProvider>
      <StyledEngineProvider injectFirst>
        <ThemeProvider theme={theme}>
          <CssBaseline />
          <Container maxWidth="lg">
            <MarginDiv>
              <h1>PPH Watchman</h1>
              <p>
                PPH Watchman is a service which downloads, parses and indexes numerous trade, government
                and non-profit lists of blocked individuals and entities to comply with those regions
                laws.
              </p>
              <p>
                <StyledLink href="https://github.com/SecurityPPH/watchman">
                  GitHub
                </StyledLink>{" "}
                |{" "}
                <StyledLink href="https://moov-io.github.io/watchman/">
                  Documentation
                </StyledLink>{" "}
                |{" "}
                <StyledLink href="https://moov-io.github.io/watchman/api/">
                  API Endpoints
                </StyledLink>
              </p>
            </MarginDiv>
            <Form onSubmit={handleSubmit} onReset={handleReset} />
            <Results data={state} />
          </Container>
        </ThemeProvider>
      </StyledEngineProvider>
    </StylesProvider>
  );
}

export default App;