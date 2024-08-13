import { Container, Grid, NavLink } from "@mantine/core";

export default function Assistants() {
  return (
    <>
      <Container size="xl">
        <Grid>
          <Grid.Col bg="blue" span={{ base: 12, md: 6, lg: 3 }}>
            <NavLink
              href="/threads.html?thread-id=1&assistant-id=1"
              label="Thread"
            ></NavLink>
          </Grid.Col>
          <Grid.Col bg="cyan" span={{ base: 12, md: 6, lg: 3 }}>
            <NavLink
              href="/threads.html?thread-id=1&assistant-id=2"
              label="Thread"
            ></NavLink>
          </Grid.Col>
          <Grid.Col bg="grape" span={{ base: 12, md: 6, lg: 3 }}>
            <NavLink
              href="/threads.html?thread-id=1&assistant-id=3"
              label="Thread"
            ></NavLink>
          </Grid.Col>
          <Grid.Col bg="indigo" span={{ base: 12, md: 6, lg: 3 }}>
            <NavLink
              href="/threads.html?thread-id=1&assistant-id=4"
              label="Thread"
            ></NavLink>
          </Grid.Col>
        </Grid>
      </Container>
    </>
  );
}
