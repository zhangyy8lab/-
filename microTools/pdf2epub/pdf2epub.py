#!/usr/bin/env python3 

import os
import argparse
from pypdf import PdfReader
from ebooklib import epub

def convert_pdf_to_epub(pdf_path):
    # 检查文件是否存在
    if not os.path.exists(pdf_path):
        print(f"Error: The file '{pdf_path}' does not exist.")
        return

    # 定义输出 EPUB 文件的路径
    output_epub = os.path.join(os.getcwd(), os.path.splitext(os.path.basename(pdf_path))[0] + '.epub')

    # 创建 EPUB 文件
    book = epub.EpubBook()
    # 设置书籍的封面和其他信息（可选）
    book.set_title(os.path.splitext(os.path.basename(pdf_path))[0])  
    book.set_language('en')
    # book.set_author("Your Author Name")  # 设置作者名
    # book.set_identifier("unique-identifier", id_type="uuid")  # 设置唯一标识符

    # 读取 PDF 文件
    reader = PdfReader(pdf_path)

    # 添加每一页作为章节
    for i, page in enumerate(reader.pages):
        # 提取文本
        text = page.extract_text()
        if text:
            # 创建 EPUB 章节
            chapter = epub.EpubHtml(title=f'Page {i + 1}', file_name=f'page_{i + 1}.xhtml', lang='en')
            chapter.set_content(f'<h1>Page {i + 1}</h1><p>{text}</p>')
            book.add_item(chapter)
            book.add_item(epub.EpubNav())

    # 设置书籍的封面和其他信息（可选）
    book.add_item(epub.EpubNcx())
    book.add_item(epub.EpubNav())

    # 保存 EPUB 文件
    epub.write_epub(output_epub, book)

    print(f"Successfully converted '{pdf_path}' to '{output_epub}'")

if __name__ == "__main__":
    # 创建命令行参数解析器
    parser = argparse.ArgumentParser(description="Convert PDF to EPUB")
    parser.add_argument("pdf_path", help="Path to the PDF file to convert")

    # 解析命令行参数
    args = parser.parse_args()

    # 执行转换
    convert_pdf_to_epub(args.pdf_path)
