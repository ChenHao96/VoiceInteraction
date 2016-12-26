public interface RecordDataListener {

    /**
     * 缓存数据的长度
     */
    int BUF_SIZE = 16 * 1024;

    /**
     * 数据处理将读取到的数据分析解析
     *
     * @param data 缓存数据
     */
    void processRecordData(byte[] data);
}
