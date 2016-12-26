public class BaseRecord4Java extends Record4Java {

    private int checkCount = 0;

    @Override
    public void processRecordData(byte[] data) {

        if (isRecordStop() || checkCount++ <= 1) return;

        int frequency = 0;
        int frequencyHave = 0;
        for (byte b : data) {
            if (Math.abs(b) < 100) {
                if (frequency++ > BUF_SIZE * 0.3) {
                    actionPerformed(Method.stop);
                    checkCount = 0;
                    return;
                }
            } else {
                if (frequencyHave++ > 400) {
                    checkCount = 0;
                    return;
                }
            }
        }
    }

    public static void main(String[] args) {

        Record4Java record4Java = new BaseRecord4Java();
        System.out.println("开始录音");
        long startTime = System.currentTimeMillis();
        record4Java.actionPerformed(Record4Java.Method.capture);
        System.out.printf("录音时长:%s毫秒\n", System.currentTimeMillis() - startTime);
        System.out.println("录音结束");
        System.out.println("播放录音");
        record4Java.actionPerformed(Record4Java.Method.play);
        System.out.println("播放结束");
    }
}
