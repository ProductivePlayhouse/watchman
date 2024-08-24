import React, { useReducer } from "react";
import * as R from "ramda";
import styled, { css } from "styled-components";
import { matchToPercent, isNilOrEmpty } from "utils";
import { Remarks } from "./Remarks";
import * as C from "Components";
import { getSDNAlts, getSDNAddresses } from "api";
import { SDNExpandDetails } from "./SDNDetails";

import Accordion from "@mui/material/Accordion";
import AccordionSummary from "@mui/material/AccordionSummary";
import AccordionDetails from "@mui/material/AccordionDetails";
import ExpandMoreIcon from "@mui/icons-material/ExpandMore";

// Styled Components
const HeaderContainer = styled.div`
  margin-top: 1em;
  width: 100%;
`;

const GridWrapper = styled.div`
  width: 100%;
  display: grid;
  grid-template-columns: 4em 1fr 1fr 1fr 4em 36px;
`;

const Row = styled.div`
  width: 100%;
  display: grid;
  padding-bottom: 1em;
  & > div {
    margin-right: 1em;
  }
`;

const SDNRow = styled.div`
  width: 100%;
  display: grid;
  grid-template-columns: 4em 1fr 1fr 1fr 4em;
  padding: 0.5em 0;
`;

const CapitalizedText = styled.div`
  text-transform: capitalize;
`;

// Styled Expansion Panel Components
const FilterExpansionPanel = ({ isExpanded, ...props }) => <Accordion {...props} />;

const ExpansionPanel = styled(FilterExpansionPanel)`
  && {
    padding: 8px;
    box-shadow: unset;
    border: 1px solid transparent;
    border-bottom: 1px solid #eee;
    ${({ isExpanded }) =>
    isExpanded &&
    `
        background-color: #f9f9f9;
        border: 1px solid #ddd;
        border-top-color: transparent;
    `};
  }
  &&:hover {
    background-color: #f9f9f9;
    border: 1px solid #ccc;
  }
  &&:last-child,
  &&:first-child {
    border-radius: 0;
  }
  &&&& * {
    user-select: initial;
  }
`;

const ExpansionPanelSummary = styled(AccordionSummary)`
  && {
    padding: 0px;
  }
`;
const ExpansionPanelDetails = styled(AccordionDetails)`
  && {
    padding: 0px;
  }
`;

const statusList = ["PRE_INIT", "INIT", "SUCCESS", "ERROR"];
const status = R.zipObj(statusList, statusList);

const initialState = {
  expanded: false,
  loaded: false,
  ALTS: {
    status: status.PRE_INIT,
    data: null,
  },
  ADDS: {
    status: status.PRE_INIT,
    data: null,
  },
};

const reducer = (state, action) =>
{
  switch (action.type)
  {
    case "EXPANDED_STATE":
      return R.assoc("expanded", action.value, state);
    case status.INIT:
      return R.assocPath([action.api, "status"], action.type, state);
    case status.SUCCESS:
      return R.pipe(
        R.assoc("loaded", true),
        R.assocPath([action.api, "status"], action.type),
        R.assocPath([action.api, "data"], action.payload || [])
      )(state);
    default:
      return state;
  }
};

export const SDN = ({ data }) =>
{
  const [details, dispatch] = useReducer(reducer, initialState);

  const handleExpandToggle = (_, expanded) =>
  {
    dispatch({ type: "EXPANDED_STATE", value: expanded });
    if (details.loaded) return;

    dispatch({ api: "ALTS", type: status.INIT });
    getSDNAlts(data.entityID).then((alts) =>
      dispatch({ api: "ALTS", type: status.SUCCESS, payload: alts })
    );

    dispatch({ api: "ADDS", type: status.INIT });
    getSDNAddresses(data.entityID).then((adds) =>
      dispatch({ api: "ADDS", type: status.SUCCESS, payload: adds })
    );
  };

  if (isNilOrEmpty(data)) return null;
  return (
    <div>
      <ExpansionPanel onChange={handleExpandToggle} isExpanded={details.expanded}>
        <ExpansionPanelSummary expandIcon={<ExpandMoreIcon />}>
          <div>
            <SDNRow>
              <div>{data.entityID}</div>
              <div>{data.sdnName}</div>
              <CapitalizedText>
                {data.sdnType || <C.Unknown>Unknown Type</C.Unknown>}
              </CapitalizedText>
              <div>{data.ofacProgram}</div>
              <div>{matchToPercent(data.match)}</div>
            </SDNRow>

            {data.sdnType === "individual" && data.title && (
              <Row style={{ paddingTop: "0.5em", color: "#666", gridTemplateColumns: "4em 1fr" }}>
                <div />
                <div>{data.title}</div>
              </Row>
            )}

            {data.sdnType === "vessel" && (
              <Row style={{ gridTemplateColumns: "4em 1fr 1fr 1fr 4em" }}>
                <div />
                <div>{data.vesselFlag || <C.Unknown>Unknown Flag</C.Unknown>}</div>
                <div>{data.vesselType || <C.Unknown>Unknown Type</C.Unknown>}</div>
                <div>{data.vesselOwner || <C.Unknown>Unknown Owner</C.Unknown>}</div>
              </Row>
            )}
            <Remarks remarks={data.remarks} />
          </div>
        </ExpansionPanelSummary>
        <ExpansionPanelDetails>
          <SDNExpandDetails data={details} />
        </ExpansionPanelDetails>
      </ExpansionPanel>
    </div>
  );
};

const Header = () => (
  <HeaderContainer>
    <GridWrapper>
      <C.ResultHeader>ID</C.ResultHeader>
      <C.ResultHeader>Name</C.ResultHeader>
      <C.ResultHeader>Type</C.ResultHeader>
      <C.ResultHeader>Program</C.ResultHeader>
      <C.ResultHeader>Match</C.ResultHeader>
      <C.ResultHeader />
    </GridWrapper>
  </HeaderContainer>
);

export const SDNS = ({ data }) =>
{
  if (!data) return null;
  return (
    <C.Section>
      <C.SectionTitle>Specially Designated Nationals</C.SectionTitle>
      <Header />
      {data && data.length > 0 && data.map((s) => <SDN key={s.entityID} data={s} />)}
    </C.Section>
  );
};
