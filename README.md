## Instructions about usage / Инструкции использования:

1. clone the given repository/ клонируйте данный репозиторий
2. go to the directory togzhan-project and start the program by executing the command: **go run main.go** / перейдите в папку togzhan-project и запустите программу, выполнив команду: **go run main.go**
3. answer question about number of go routines, the program will repeat the question until valid answer is provided / ответьте на вопрос о количестве горутин в консоле, программа будет спрашивать вас вопрос до тех пор, пока вы не дадите валидный ответ (число > 0)
4. then answer question about regeneration of file, write either _y_ or _n_, the answer is case insensitive / после этого ответьте на вопрос о пересоздании файла, напишите _y (если да)_ или _n (если нет)_, ответ может быть, как и с заглавной, так и с маленькой буквы
   * if you did not want to regenerate file, expect that pregenerated _sample_file_1m.json_ with **1,000,000** objects will be used / если вы не захотели пересоздавать файл, ожидайте, что программа будет использовать файл уже созданный _sample_file_1m.json_ с **1,000,000** объектами
   * if you chose to regenerate file, then expect that _sample_file_1m.json_ file with **1,000,000** objects will be re-created and used. I would chose this option because you might have problems with accessing pre-generated file because of its size / если вы выбрали пересоздать файл, ожидайте, что файл _sample_file_1m.json_ с **1,000,000** объектами будет создан заново и использоваться для вычисления программой. Я бы выбрала вариант пересоздания, чтобы избежать проблем доступа к существующему файлу из-за его размера
5. done, you should see sum in you terminal / все, вы должны увидеть сумму чисел в вашей консоли

 All the comments in the code is written in English, because it's the suitable for author work language (apart from go itself:)
 Все комментарии в коде написаны на английском языке, потому что это удобный для автора язык работы (не считая самого Go:)
