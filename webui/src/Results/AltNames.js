import React from "react";
import styled from "styled-components";
import { matchToPercent } from "../utils";
import { Remarks } from "./Remarks";
import * as C from "Components";

// Styled Components
const HeaderWrapper = styled.div`
  margin-top: 1em;
  width: 100%;
  display: grid;
  grid-template-columns: ${({ withMatch }) =>
    withMatch ? `4em 3fr 1fr 4em` : `4em 4fr 1fr`};
`;

const AltNameWrapper = styled.div`
  padding-bottom: 1em;
  border-bottom: 1px solid #eee;
  &:last-of-type {
    border-bottom: 0px;
  }
  & > div {
    border: 0px solid red;
  }
`;

const AltNameGrid = styled.div`
  display: grid;
  grid-template-columns: ${({ hasMatch }) =>
    hasMatch ? `4em 3fr 1fr 4em` : `4em 4fr 1fr`};
  padding-top: 1em;
  & > div {
    margin-right: 1em;
  }
`;

const UppercaseText = styled.div`
  text-transform: uppercase;
`;

// Components
export const Header = ({ withMatch = true }) => (
  <HeaderWrapper withMatch={withMatch}>
    <C.ResultHeader>ID</C.ResultHeader>
    <C.ResultHeader>Alternate Name</C.ResultHeader>
    <C.ResultHeader>Type</C.ResultHeader>
    {withMatch && <C.ResultHeader>Match</C.ResultHeader>}
  </HeaderWrapper>
);

export const AltNames = ({ data }) =>
{
  if (!data) return null;
  return (
    <C.Section>
      <C.SectionTitle>Alternate Names</C.SectionTitle>
      <Header />
      {data.length > 0 &&
        data.map((s) => <AltName key={`${s.entityID}-${s.alternateID}`} data={s} />)}
    </C.Section>
  );
};

export const AltName = ({ data, displayId = "entityID" }) =>
{
  return (
    <AltNameWrapper>
      <AltNameGrid hasMatch={!!data.match}>
        <div>{data[displayId]}</div>
        <div>{data.alternateName}</div>
        <UppercaseText>{data.alternateType}</UppercaseText>
        {data.match && <div>{matchToPercent(data.match)}</div>}
      </AltNameGrid>
      <Remarks remarks={data.remarks} />
    </AltNameWrapper>
  );
};
