import google.generativeai as genai

genai.configure(api_key="AIzaSyBlECYIGkDFbuDATCml7OM179ioNezt6uM")

#print("Modelos disponíveis para você:")
#for m in genai.list_models():
#    if 'generateContent' in m.supported_generation_methods:
#        print(m.name)

model = genai.GenerativeModel('models/gemini-2.5-pro')

try:
    response = model.generate_content("Olá! Testando conexão.")
    print(response.text)
except Exception as e:
    print(f"Erro ao gerar conteúdo: {e}")


#models/gemini-2.5-flash
#models/gemini-2.5-pro
#models/gemini-2.0-flash
#models/gemini-2.0-flash-001
#models/gemini-2.0-flash-lite-001
#models/gemini-2.0-flash-lite
#models/gemini-2.5-flash-preview-tts
#models/gemini-2.5-pro-preview-tts
#models/gemma-3-1b-it
#models/gemma-3-4b-it
#models/gemma-3-12b-it
#models/gemma-3-27b-it
#models/gemma-3n-e4b-it
#models/gemma-3n-e2b-it
#models/gemini-flash-latest
#models/gemini-flash-lite-latest
#models/gemini-pro-latest
#models/gemini-2.5-flash-lite
#models/gemini-2.5-flash-image
#models/gemini-2.5-flash-lite-preview-09-2025
#models/gemini-3-pro-preview
#models/gemini-3-flash-preview
#models/gemini-3.1-pro-preview
#models/gemini-3.1-pro-preview-customtools
#models/gemini-3.1-flash-lite-preview
#models/gemini-3-pro-image-preview
#models/nano-banana-pro-preview
#models/gemini-3.1-flash-image-preview
#models/lyria-3-clip-preview
#models/lyria-3-pro-preview
#models/gemini-robotics-er-1.5-preview
#models/gemini-2.5-computer-use-preview-10-2025
#models/deep-research-pro-preview-12-2025