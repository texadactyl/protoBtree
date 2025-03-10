import java.io.*;
import java.util.*;

class IndexRecord {
    int key;  // Record number
    long value; // Byte offset

    public IndexRecord(int key, long value) {
        this.key = key;
        this.value = value;
    }
}

public class reader {
    private static final String INDEX_FILE_PATH = "index.bin";
    private static final String DATA_FILE_PATH = "data.bin";

    public static void main(String[] args) {
        try {
            // Load B-tree index as raw binary
            Map<Integer, Long> bTree = loadBTreeAsBinary(INDEX_FILE_PATH);
            System.out.println("Loaded B-tree: " + bTree);

            // Read records using record numbers 42, 3, 17
            int[] recordNumbers = {42, 3, 17};
            for (int recordNumber : recordNumbers) {
                readDataByRecordNumber(DATA_FILE_PATH, bTree, recordNumber);
            }
        } catch (Exception e) {
            e.printStackTrace();
        }
    }

    // Read index file as raw binary
    private static Map<Integer, Long> loadBTreeAsBinary(String filePath) throws IOException {
        Map<Integer, Long> bTree = new HashMap<>();
        File file = new File(filePath);
        if (!file.exists()) {
            System.out.println("Index file not found.");
            return bTree;
        }

        try (DataInputStream dis = new DataInputStream(new FileInputStream(file))) {
            while (dis.available() >= 12) { // Each record is 4 + 8 bytes
                int key = dis.readInt();   // Read int (4 bytes)
                long value = dis.readLong(); // Read long (8 bytes)
                bTree.put(key, value);
            }
        }
        return bTree;
    }

    // Read a record from the data file using the index
    private static void readDataByRecordNumber(String dataPath, Map<Integer, Long> bTree, int recordNumber) {
        Long offset = bTree.get(recordNumber);
        if (offset == null) {
            System.out.println("Record number " + recordNumber + " not found in the index.");
            return;
        }

        try (RandomAccessFile dataFile = new RandomAccessFile(dataPath, "r")) {
            dataFile.seek(offset); // Move to the record's offset

            // Read and print the record as a string (assuming text data)
            byte[] buffer = new byte[100]; // Adjust size as needed
            int bytesRead = dataFile.read(buffer);
            if (bytesRead > 0) {
                System.out.println("Record " + recordNumber + ": " + new String(buffer, 0, bytesRead));
            } else {
                System.out.println("Record " + recordNumber + " is empty.");
            }
        } catch (IOException e) {
            System.out.println("Error reading record at offset " + offset + ": " + e.getMessage());
        }
    }
}

