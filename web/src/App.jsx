import {
  Button,
  Center,
  Flex,
  MantineProvider,
  createTheme,
} from "@mantine/core";
import "@mantine/core/styles.css";

export default function App() {
  const theme = createTheme({});
  return (
    <>
      <MantineProvider theme={theme} defaultColorScheme="auto">
        <Center m="xl">
          <Flex
            mih={50}
            bg="rgba(0, 0, 0, .3)"
            gap="md"
            justify="flex-start"
            align="flex-start"
            direction="row"
            wrap="wrap"
          >
            <h1>Hello </h1>
            <Button color="red">Button 1</Button>
            <Button color="indigo.5">Button 2</Button>
            <Button color="violet">Button 3</Button>
          </Flex>
        </Center>
      </MantineProvider>
    </>
  );
}
