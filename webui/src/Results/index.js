import React from "react";
import styled from "styled-components";
import * as C from "../Components";
import { HighestMatch } from "./HighestMatch";
import { SDNS } from "./SDN";
import { AltNames } from "./AltNames";
import { Addresses } from "./Addresses";
import { DeniedPersons } from "./DeniedPersons";
import { isNilOrEmpty } from "utils";
import { SectoralSanctions } from "./SectoralSanctions";

// Styled Components
const CenteredContainer = styled.div`
  display: flex;
  justify-content: center;
`;

const render = ({ data }) =>
{
  const { loading, error, results } = data;
  if (loading)
    return (
      <C.Container>
        <CenteredContainer>
          <C.Spinner />
        </CenteredContainer>
      </C.Container>
    );
  if (error) return <C.Container>ERROR: {error.message}</C.Container>;
  if (isNilOrEmpty(results)) return null;
  return (
    <C.Container>
      <HighestMatch data={results.highestMatch} />
      <SDNS data={results.SDNs} />
      <AltNames data={results.altNames} />
      <Addresses data={results.addresses} />
      <DeniedPersons data={results.deniedPersons} />
      <SectoralSanctions data={results.sectoralSanctions} />
    </C.Container>
  );
};

export default render;
