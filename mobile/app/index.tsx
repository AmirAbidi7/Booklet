import { Text, View } from "react-native";
import Pdf from "react-native-pdf";

export default function Index() {
  return (
    <View className="flex-1 justify-center items-center">
      <Text className="text-xl">Edit app/index.tsx to edit this screen.</Text>
      <Pdf
        source={{
          uri: "https://assets.withfra.me/pdf/sample.pdf",
          cache: true,
        }}
      />
    </View>
  );
}
