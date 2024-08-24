import React from "react";
import * as R from "ramda";
import styled from "styled-components";
import CircularProgress from "@mui/material/CircularProgress";
import * as C from "Components";
import { Header as AddressHeader, Address } from "./Addresses";
import { Header as AlternatesHeader, AltName } from "./AltNames";

// Styled Components
const SectionContainer = styled.div`
  margin-bottom: 1em;
`;

const SDNExpandDetailsContainer = styled.div`
  width: 100%;
  & > * {
    margin-bottom: 1.5em;
  }
  & > *:last-child {
    margin-bottom: 0;
  }
`;

const Addresses = ({ data }) =>
{
  if (!data) return <CircularProgress size="1em" />;
  if (data.length === 0) return null;
  return (
    <SectionContainer>
      <C.SectionTitle>Addresses</C.SectionTitle>
      <AddressHeader withMatch={false} />
      {data.map((a) => (
        <Address key={a.addressID} data={a} displayId="addressID" />
      ))}
    </SectionContainer>
  );
};

const Alternates = ({ data }) =>
{
  if (!data) return <CircularProgress size="1em" />;
  if (data.length === 0) return null;
  return (
    <SectionContainer>
      <C.SectionTitle>Alternate Names</C.SectionTitle>
      <AlternatesHeader withMatch={false} />
      {data.map((a) => (
        <AltName key={a.alternateID} data={a} displayId="alternateID" />
      ))}
    </SectionContainer>
  );
};

export const SDNExpandDetails = ({ data }) => (
  <SDNExpandDetailsContainer>
    <Addresses data={R.path(["ADDS", "data"])(data)} />
    <Alternates data={R.path(["ALTS", "data"])(data)} />
  </SDNExpandDetailsContainer>
);
