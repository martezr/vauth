import * as React from 'react';
import { 
  Grid,
  GridItem,
  Page, 
  PageSection,
  PageSectionVariants,
  Text,
  TextContent,
  Title } from '@patternfly/react-core';
  import { Card, CardTitle, CardBody, CardFooter, Gallery, GalleryItem } from '@patternfly/react-core';
  import { ChartDonutThreshold, ChartDonutUtilization } from '@patternfly/react-charts';

const Dashboard: React.FunctionComponent = () => (
  <Page>
  <PageSection variant={PageSectionVariants.light}>
  <TextContent>
              <Text component="h1">Dashboard</Text>
            </TextContent>
  </PageSection>
    <PageSection>
            <Grid hasGutter>
    <GridItem span={8}>
    <Gallery hasGutter minWidths={{ default: '360px' }}>
    <GalleryItem>
      <Card id="utilization-card-1" component="div">
        <CardTitle>
          <Title headingLevel="h2" size="lg">
            CPU Usage
          </Title>
        </CardTitle>
        <CardBody>
          <ChartDonutThreshold
            ariaDesc="Mock storage capacity"
            ariaTitle="Mock donut utilization chart"
            constrainToVisibleArea={true}
            data={[
              { x: 'Warning at 60%', y: 60 },
              { x: 'Danger at 90%', y: 90 }
            ]}
            height={200}
            labels={({ datum }) => (datum.x ? datum.x : null)}
            padding={{
              bottom: 0,
              left: 10,
              right: 150,
              top: 0
            }}
            width={350}
          >
            <ChartDonutUtilization
              data={{ x: 'Storage capacity', y: 80 }}
              labels={({ datum }) => (datum.x ? `${datum.x}: ${datum.y}%` : null)}
              legendData={[{ name: `Capacity: 80%` }, { name: 'Warning at 60%' }, { name: 'Danger at 90%' }]}
              legendOrientation="vertical"
              title="80%"
              subTitle="of 100 GBps"
              thresholds={[{ value: 60 }, { value: 90 }]}
            />
          </ChartDonutThreshold>{' '}
        </CardBody>
        <CardFooter>
          <a href="#">See details</a>
        </CardFooter>
      </Card>
    </GalleryItem>
  </Gallery>

    </GridItem>
    <GridItem span={4} rowSpan={2}>
      span = 4, rowSpan = 2
    </GridItem>
    <GridItem span={2} rowSpan={3}>
      span = 2, rowSpan = 3
    </GridItem>
    <GridItem span={2}>span = 2</GridItem>
    <GridItem span={4}>span = 4</GridItem>
    <GridItem span={2}>span = 2</GridItem>
    <GridItem span={2}>span = 2</GridItem>
    <GridItem span={2}>span = 2</GridItem>
    <GridItem span={4}>span = 4</GridItem>
    <GridItem span={2}>span = 2</GridItem>
    <GridItem span={4}>span = 4</GridItem>
    <GridItem span={4}>span = 4</GridItem>
  </Grid>

            </PageSection>
  </Page>
)

export { Dashboard };
