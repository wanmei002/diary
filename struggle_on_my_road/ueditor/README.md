### go ʵ�� �ٶȸ��ı���(ueditor)
> û��ʲô�ѵ� ��Ҫ�ο� ueditor PHP ��

> ץȡ ueditor php�� ����ĺ�̨���� �� golang �滻

1. ueditor.config.js ����ĺ�̨��ַ�滻��
    + �� `serverUrl: URL + "php/controller.php"` �滻�� `serverUrl: /controller` ,ueditor һ�㶼�Ǹ������ַ����
    + url ���� action ��ֵ������̨���� ����Ҫʵ���� ���������ļ�  �ϴ�ͼƬ�Ĺ���
    + �Ұ� /controller ��ַ������ UEditor.ControllerUE �������������Ҫʵ���� ����������Ϣ �� �ϴ�ͼƬ����~~~ԭ���������
    + ����ʵ���뿴��Ŀ¼�µ� `UEditor.go` �ļ�
    
### �ļ�����
 + `Create.html` ���ı���� html ҳ��
 + `config.json` ��˷��ص� json �ļ�
 + `ueditor.config.js` ueditor ��js �����ļ�
 + `UEditor.go` ����ļ�, ��Ҫʵ���� ���� ������Ϣ �� �ϴ�ͼƬ�Ĺ���